import 'dart:developer';

import 'package:mongo_dart/mongo_dart.dart';
import 'package:trashmap/models/consts.dart';

class MongoDBService {
  static var _db, _collection;

  static connect() async {
    _db = await Db.create(MONGODB_URL);
    await _db.open();
    print(_db);
    _collection = _db.collection(COLL_NAME);
  }

  // Future<List<Map<String, dynamic>>> fetchData(String collectionName) async {
  //   DbCollection collection = _db.collection(collectionName);
  //   List<Map<String, dynamic>> data = await collection.find().toList();
  //   return data;
  // }

  // void close() {
  //   _db.close();
  // }
}
