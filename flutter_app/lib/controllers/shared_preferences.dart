import 'dart:async';

import 'package:shared_preferences/shared_preferences.dart';


Future<SharedPrefs> generateSharedPrefs() async {
  SharedPrefs _prefs = new SharedPrefs();
  await _prefs.init();
  return _prefs;
}

class SharedPrefs {
  SharedPreferences _prefs;

  // FlutterSecureStorage _secureStorage;

  Future _doneFuture;

  SharedPrefs();

  Future<void> init() async {
    _prefs = await SharedPreferences.getInstance();
    // _secureStorage = new FlutterSecureStorage();
    return;
  }

  static const QUICK_WORD = 'quick-word';
  static const TOTAL_STUDY_TIME = 'total-study-time';
  static const REGISTRATION_DATE = 'registration-date';
  static const CURRENT_ROOM_ID = 'current-room-id';
  static const CURRENT_ROOM_NAME = 'current-room-name';


  Future<bool> setQuickWord(String quickWord) {
    return _prefs.setString(QUICK_WORD, quickWord);
  }
  Future<String> getQuickWord() async {
    return _prefs.getString(QUICK_WORD) ?? '';
  }

  Future<bool> setTotalStudyTime(int sumStudyTime) {
    return _prefs.setInt(TOTAL_STUDY_TIME, sumStudyTime);
  }
  int getTotalStudyTime() {
    return _prefs.getInt(TOTAL_STUDY_TIME) ?? 0;
  }

  Future<bool> setRegistrationDate(DateTime date) {
    return _prefs.setInt(REGISTRATION_DATE, date.millisecondsSinceEpoch);
  }
  DateTime getRegistrationDate() {
    return DateTime.fromMillisecondsSinceEpoch(_prefs.getInt(REGISTRATION_DATE) ?? 0);
  }

  Future<bool> setCurrentRoomId(String roomId) {
    return _prefs.setString(CURRENT_ROOM_ID, roomId);
  }
  Future<String> getCurrentRoomId() async {
    return _prefs.getString(CURRENT_ROOM_ID) ?? '';
  }

  Future<bool> setCurrentRoomName(String roomName) {
    return _prefs.setString(CURRENT_ROOM_NAME, roomName);
  }
  Future<String> getCurrentRoomName() async {
    return _prefs.getString(CURRENT_ROOM_NAME) ?? '';
  }
}