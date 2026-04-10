import 'dart:async';
import 'dart:typed_data';
import 'package:web_socket_channel/web_socket_channel.dart';
import 'package:web_socket_channel/status.dart' as status;

class WebSocket {
  WebSocketChannel? _channel;
  StreamSubscription? _subscription;

  Function(Uint8List binaryData)? onMessageReceived;
  Function()? onConnected;
  Function(String reason)? onDisconnected;
  Function(String error)? onError;

  final String token;
  final String baseUrl;

  WebSocket({required this.token, required this.baseUrl});

  Future<void> connect() async {
    try {
      final wsUrl = Uri.parse('$baseUrl/ws/$token');
      _channel = WebSocketChannel.connect(wsUrl);

      await _channel!.ready;

      onConnected?.call();

      // Listen to incoming messages
      _subscription = _channel!.stream.listen(
        (message) {
          if (message is Uint8List) {
            onMessageReceived?.call(message);
          } else if (message is String) {
            print('Text message from server: $message');
          }
        },
        onError: (error) {
          onError?.call(error.toString());
          _handleDisconnect();
        },
        onDone: () {
          _handleDisconnect();
        },
      );
    } catch (e) {
      onError?.call('Connection failed: $e');
    }
  }

  void sendInput(Uint8List inputData) {
    _channel?.sink.add(inputData);
  }

  void disconnect() {
    _channel?.sink.close(status.normalClosure, 'User disconnected');
    _subscription?.cancel();
    _channel = null;
  }

  void _handleDisconnect() {
    onDisconnected?.call('Connection closed');
    _subscription?.cancel();
    _channel = null;
  }
}
