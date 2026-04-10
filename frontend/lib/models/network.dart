import 'package:frontend/utils/network_request.dart';

class Network {
  final String baseUrl;

  Network({required this.baseUrl});

  Future<RespServerInfo?> loadStatus() async {
    NetworkRequest request = NetworkRequest<Map<String, dynamic>>();
    await request.fetch("status", baseUrl: baseUrl);
    if (request.hasError) {
      return null;
    }

    try {
      return RespServerInfo.fromJson(request.data!);
    } catch (e) {
      return null;
    }
  }

  Future<RespAuth?> getAuthToken(String name, int colorIndex) async {
    NetworkRequest request = NetworkRequest<Map<String, dynamic>>();
    final reqName = name.replaceAll(RegExp(r' '), ":");
    await request.fetch('auth/$reqName:$colorIndex', baseUrl: baseUrl);
    if (request.hasError) {
      return null;
    }

    try {
      return RespAuth.fromJson(request.data!);
    } catch (e) {
      return null;
    }
  }
}

class RespServerInfo {
  final int serverRunningMin;
  final List<String> nameAdvList;
  final List<String> nameList;

  RespServerInfo({
    this.serverRunningMin = 0,
    this.nameAdvList = const [],
    this.nameList = const [],
  });

  static RespServerInfo fromJson(Map<String, dynamic> json) {
    return RespServerInfo(
      serverRunningMin: json["serverRunningMin"],
      nameAdvList: List<String>.from(json["advList"]),
      nameList: List<String>.from(json["nameList"]),
    );
  }
}

class RespAuth {
  final String token;
  RespAuth({required this.token});

  static RespAuth fromJson(Map<String, dynamic> json) {
    return RespAuth(token: json["token"]);
  }
}
