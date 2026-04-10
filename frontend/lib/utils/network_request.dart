import 'dart:convert';

import 'package:http/http.dart' as http;

class NetworkRequest<T> {
  bool hasError = false;
  String errorMessage = '';
  int statusCode = 0;
  T? data;

  Future<void> fetch(
    String endpoint, {
    String baseUrl = 'http://localhost:8080',
  }) async {
    hasError = false;
    data = null;
    errorMessage = '';

    try {
      final response = await http.get(Uri.parse('$baseUrl/$endpoint'));

      statusCode = response.statusCode;

      if (response.statusCode == 200) {
        data = jsonDecode(response.body) as T?;
      } else {
        hasError = true;
        errorMessage = 'Server error: ${response.statusCode}';
      }
    } catch (e) {
      hasError = true;
      errorMessage = 'Connection error: $e';
    }
  }
}
