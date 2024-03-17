import 'dart:async';
// import 'dart:html';

import 'package:flutter/material.dart';
import 'package:google_maps_flutter/google_maps_flutter.dart';
import 'package:location/location.dart';
import 'package:trashmap/services/mongoDB_service.dart';
import 'package:trashmap/models/trashbin.dart';

class MapPage extends StatefulWidget {
  const MapPage({Key? key}) : super(key: key);

  @override
  State<MapPage> createState() => _MapPageState();
}

class _MapPageState extends State<MapPage> {
  final Location _locationController = Location();

  final Completer<GoogleMapController> _mapController =
      Completer<GoogleMapController>();

  static const LatLng _pLeipzig = LatLng(51.3418814, 12.3731);
  LatLng? _currentPos = null;

  @override
  void initState() {
    super.initState();
    getLocationUpdates();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: _currentPos == null
          ? const Center(child: Text("Loading ..."))
          : FutureBuilder<Set<Marker>>(
              future: _createMarkers(),
              builder:
                  (BuildContext context, AsyncSnapshot<Set<Marker>> snapshot) {
                if (snapshot.hasData) {
                  return GoogleMap(
                    onMapCreated: (GoogleMapController controller) {
                      _mapController.complete(controller);
                    }, // mapCreated passes the controller of the corresponding map. I then store it in the _mapController variable

                    initialCameraPosition: const CameraPosition(
                      target: _pLeipzig,
                      zoom: 13.0,
                    ),
                    markers: snapshot.data!,
                  );
                } else if (snapshot.hasError) {
                  return Text('Error: ${snapshot.error}');
                } else {
                  return const CircularProgressIndicator();
                }
              },
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

  Future<List<Trashbin>> _getTrashbins() async {
    return await DBService().getAllTrashbins();
  }

  Future<Set<Marker>> _createMarkers() async {
    List<Trashbin> trashbins = await _getTrashbins();
    Set<Marker> markers = trashbins.map((trashbin) {
      return Marker(
        markerId: MarkerId(trashbin.id),
        position: LatLng(trashbin.lat, trashbin.lng),
        icon: BitmapDescriptor.defaultMarker,
      );
    }).toSet();

    markers.add(
      Marker(
        markerId: const MarkerId("_currentLocation"),
        icon: BitmapDescriptor.defaultMarker,
        position: _currentPos!,
      ),
    );

    // markers.addAll(
    //   latlngList.map((latlng) {
    //     return Marker(
    //       markerId: MarkerId('${latlng['lat']},${latlng['lng']}'),
    //       position: LatLng(latlng['lat']!, latlng['lng']!),
    //       icon: BitmapDescriptor.defaultMarker,
    //     );
    //   }).toSet(),
    // );

    return markers;
  }
}
