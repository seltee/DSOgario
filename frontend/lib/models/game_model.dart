import 'dart:math';
import 'dart:typed_data';

import 'package:flutter/widgets.dart';
import 'package:frontend/core/blob_colors.dart';
import 'package:frontend/models/config.dart';
import 'package:frontend/models/game_world.dart';
import 'package:frontend/models/network.dart';
import 'package:frontend/models/theme_model.dart';
import 'package:frontend/utils/provider.dart';
import 'package:frontend/utils/web_socket.dart';

final random = Random();

const baseUrl = 'http://localhost:8080';
const baseWSUrl = 'ws://localhost:8080';

enum GameState { loading, auth, game }

class GameModel extends ChangeNotifier {
  Network network = Network(baseUrl: baseUrl);
  GameWorld gameWorld = GameWorld();
  ThemeModel themeModel = ThemeModel();

  RespServerInfo serverInfo = RespServerInfo();

  bool gameReady = false;
  bool isLoggingIn = false;
  bool isInGame = false;
  bool serverIsUnavailable = false;
  String name = "";
  String token = "";

  int colorIndex = 0;

  WebSocket? _webSocket;

  GameModel() {
    gameWorld.sendPlayerDirection = (Offset playerDirection) =>
        sendPlayerDirection(playerDirection);
  }

  void initGame() async {
    RespServerInfo? status = await network.loadStatus();
    if (status != null) {
      serverInfo = status;
      suggestNewName();
      gameReady = true;
    } else {
      serverIsUnavailable = true;
    }
    notifyListeners();
  }

  GameState getState() {
    if (gameReady) {
      if (isInGame) {
        return GameState.game;
      } else {
        return GameState.auth;
      }
    }
    return GameState.loading;
  }

  void suggestNewName() {
    if (random.nextInt(10) == 0) {
      name = serverInfo.nameList[random.nextInt(serverInfo.nameList.length)];
    } else {
      var advName =
          serverInfo.nameAdvList[random.nextInt(serverInfo.nameAdvList.length)];
      var realName =
          serverInfo.nameList[random.nextInt(serverInfo.nameList.length)];
      name = '$advName $realName';
    }
    colorIndex = random.nextInt(blobColors.length);
    notifyListeners();
  }

  void enterTheGame() async {
    isInGame = false;
    isLoggingIn = true;
    notifyListeners();

    var newToken = (await network.getAuthToken(name, colorIndex))?.token;

    if (newToken != null) {
      token = newToken;
      _webSocket = WebSocket(token: token, baseUrl: baseWSUrl);

      _webSocket?.onConnected = () {
        isLoggingIn = false;
        isInGame = true;
        gameWorld.reset();
        notifyListeners();
      };
      _webSocket?.onMessageReceived = (Uint8List data) {
        if (data.length < 8) {
          print('Corrupted frame received');
          return;
        }
        final byteData = data.buffer.asByteData();
        int messageType = byteData.getUint16(0, Endian.big);

        switch (messageType) {
          case 1:
            gameWorld.provideFrame(byteData);
            break;
          case 2:
            gameWorld.provideScore(byteData);
            break;
          default:
            print('Unknown message type $messageType');
        }
      };
      _webSocket?.onError = (String error) {
        isLoggingIn = false;
        isInGame = false;
        gameReady = false;
        initGame();
        notifyListeners();
        print("ERROR");
        print(error);
      };
      _webSocket?.onDisconnected = (String error) {
        isLoggingIn = false;
        isInGame = false;
        gameReady = false;
        initGame();
        notifyListeners();
      };

      _webSocket?.connect();
      notifyListeners();
    } else {
      isLoggingIn = false;
      notifyListeners();
    }
  }

  void restart() {
    disconnect();
    enterTheGame();
    notifyListeners();
  }

  void goToMainMenu() {
    disconnect();
    isLoggingIn = false;
    isInGame = false;
    notifyListeners();
  }

  void disconnect() {
    _webSocket?.disconnect();
    _webSocket = null;
  }

  void update(double delta) {
    gameWorld.update(delta);
  }

  void sendPlayerDirection(Offset playerDirection) {
    final int directionX = (playerDirection.dx * precision)
        .clamp(-32768, 32767)
        .toInt();
    final int directionY = (playerDirection.dy * precision)
        .clamp(-32768, 32767)
        .toInt();

    final inputBytes = Uint8List(6);
    final ByteData bd = inputBytes.buffer.asByteData();
    bd.setUint16(0, 0, Endian.big); // message type, position
    bd.setInt16(2, directionX, Endian.big); // direction x
    bd.setInt16(4, directionY, Endian.big); // direction y

    _webSocket?.sendInput(inputBytes);
  }
}

class GameModelProvider extends Provider<GameModel> {
  const GameModelProvider({
    super.key,
    required super.child,
    required super.model,
  });

  // Easy access via context
  static GameModel of(BuildContext context) {
    final provider = context
        .dependOnInheritedWidgetOfExactType<GameModelProvider>();
    assert(provider != null, 'No ModelProvider found in context');
    return provider!.model;
  }
}
