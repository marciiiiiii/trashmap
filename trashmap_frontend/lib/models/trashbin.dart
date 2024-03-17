class Trashbin {
  final String id;
  final double lat;
  final double lng;

  Trashbin({required this.id, required this.lat, required this.lng});

  factory Trashbin.fromJson(Map<String, dynamic> json) {
    return Trashbin(
      id: json['_id'],
      lat: json['lat'],
      lng: json['lng'],
    );
  }
}
