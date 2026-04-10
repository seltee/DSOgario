import 'package:flutter/scheduler.dart';
import 'package:frontend/models/game_model.dart';

class GameController {
  final GameModel model;
  late final Ticker _ticker;
  Duration _last = Duration.zero;

  GameController(this.model, TickerProvider vsync) {
    _ticker = vsync.createTicker(_tick);
  }

  void start() => _ticker.start();
  void stop() => _ticker.stop();

  void dispose() => _ticker.dispose();

  void _tick(Duration elapsed) {
    final delta = (elapsed - _last).inMicroseconds / 1e6;
    _last = elapsed;
    model.update(delta);
  }
}
