import 'dart:convert';
import 'dart:developer';

import 'package:trashmap/models/consts.dart';
import 'package:trashmap/models/trashbin.dart';
import 'package:http/http.dart' as http;

class DBService {
  Future<List<Trashbin>> getAllTrashbins() async {
    final response = await http.get(TRASHBIN_API_URL);

    if (response.statusCode == 200) {
      List<dynamic> trashbins = json.decode(response.body);
      return trashbins.map((item) => Trashbin.fromJson(item)).toList();
    } else {
      throw Exception('Failed to load trashbins');
    }
  }
}
