import 'dart:async';
// import 'dart:html';

import 'package:flutter/material.dart';
import 'package:google_maps_flutter/google_maps_flutter.dart';
import 'package:location/location.dart';
// import 'package:flutter_polyline_points/flutter_polyline_points.dart';
// import 'package:trashmap/consts.dart';
// import 'package:http/http.dart' as http;

class MapPage extends StatefulWidget {
  const MapPage({Key? key}) : super(key: key);

  @override
  State<MapPage> createState() => _MapPageState();
}

class _MapPageState extends State<MapPage> {
  Location _locationController = Location();

  final Completer<GoogleMapController> _mapController =
      Completer<GoogleMapController>();

  static const LatLng _pLeipzig = LatLng(51.3418814, 12.3731);
  static const LatLng _pKrostitz = LatLng(51.4622649, 12.4443108);
  static const List<Map<String, double>> latlngList = [
    {"lat": 51.3418814, "lng": 12.3731},
    {"lat": 51.4622649, "lng": 12.4443108}
  ];
  LatLng? _currentPos = null;

  @override
  void initState() {
    super.initState();
    getLocationUpdates();
    // getLocationUpdates().then(
    //   (_) => {
    //     getPolylinePoints().then(
    //       (coordinates) => print(coordinates),
    //     )
    //   },
    // );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: _currentPos == null
          ? const Center(child: Text("Loading ..."))
          : GoogleMap(
              onMapCreated: (GoogleMapController controller) {
                _mapController.complete(controller);
              }, // mapCreated passes the controller of the corresponding map. I then store it in the _mapController variable

              initialCameraPosition: const CameraPosition(
                target: _pLeipzig,
                zoom: 13.0,
              ),
              markers: _createMarkers(),
            ),
    );
  }

  Future<void> _cameraToPosition(LatLng position) async {
    final GoogleMapController controller = await _mapController.future;
    controller.animateCamera(CameraUpdate.newCameraPosition(
      CameraPosition(
        target: position,
        zoom: 13.0,
      ),
    ));
  }

  Future<void> getLocationUpdates() async {
    bool serviceEnabled;
    PermissionStatus permissionGranted;

    serviceEnabled = await _locationController.serviceEnabled();

    if (serviceEnabled) {
      _locationController.requestService();
    } else {
      return;
    }

    permissionGranted = await _locationController.hasPermission();

    if (permissionGranted == PermissionStatus.denied) {
      permissionGranted = await _locationController.requestPermission();
      if (permissionGranted != PermissionStatus.granted) {
        return;
      }
    }

    _locationController.onLocationChanged
        .listen((LocationData currentLocation) {
      if (currentLocation.latitude != null &&
          currentLocation.longitude != null) {
        setState(() {
          _currentPos = LatLng(
              currentLocation.latitude!,
              currentLocation
                  .longitude!); // "!" is telling Dart that value will not be null
          _cameraToPosition(_currentPos!);
        });
      }
    });
  }

  // Future<List<LatLng>> getPolylinePoints() async {
  //   PolylinePoints polylinePoints = PolylinePoints();
  //   List<LatLng> polylineCoordinates = [];
  //   final response = await http.get(Uri.parse(
  //       'https://maps.googleapis.com/maps/api/directions/json?origin=51.3418814%2C12.3731&destination=51.4622649%2C12.4443108&mode=walking&avoidHighways=false&avoidFerries=true&avoidTolls=false&alternatives=false&key=${GOOGLE_MAPS_API_KEY}'));
  //       if (response.statusCode == 200) {
  //         PolylineResult result = PolylineResult.fromJson(jsonDecode(response.body));
  //       } else {
  //         print(response.statusCode);
  //       }
  // await polylinePoints.getRouteBetweenCoordinates(
  //   GOOGLE_MAPS_API_KEY,
  //   PointLatLng(_pLeipzig.latitude, _pLeipzig.longitude),
  //   PointLatLng(_pKrostitz.latitude, _pKrostitz.longitude),
  //   travelMode: TravelMode.walking,
  // );

  //   if (result.points.isNotEmpty) {
  //     result.points.forEach((PointLatLng point) {
  //       polylineCoordinates.add(LatLng(point.latitude, point.longitude));
  //     });
  //   } else {
  //     print(result.errorMessage);
  //   }

  //   return polylineCoordinates;
  // }

  Set<Marker> _createMarkers() {
    Set<Marker> markers = {};

    markers.add(
      Marker(
        markerId: const MarkerId("_currentLocation"),
        icon: BitmapDescriptor.defaultMarker,
        position: _currentPos!,
      ),
    );

    markers.addAll(
      latlngList.map((latlng) {
        return Marker(
          markerId: MarkerId('${latlng['lat']},${latlng['lng']}'),
          position: LatLng(latlng['lat']!, latlng['lng']!),
          icon: BitmapDescriptor.defaultMarker,
        );
      }).toSet(),
    );

    return markers;
  }
}
