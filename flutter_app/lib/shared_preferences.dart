import 'dart:async';

import 'package:shared_preferences/shared_preferences.dart';

class SharedPrefs {
  SharedPreferences _prefs;

  // FlutterSecureStorage _secureStorage;

  Future _doneFuture;

  SharedPrefs();

  Future<void> init() async {
    _prefs = await SharedPreferences.getInstance();
    // _secureStorage = new FlutterSecureStorage();
  }

  static const USES_ID = 'user-id';
  static const DISPLAY_NAME = 'display-name';
  static const QUICK_WORD = 'quick-word';
  static const ACCOUNT_TYPE = 'account-type';
  static const MAIL_ADDRESS = 'mail-address';
  static const SUM_STUDY_TIME = 'sum-study-time';
  static const REGISTRATION_DATE = 'registration-date';

  // static const ID_TOKEN = 'id-token';

  Future<bool> setUserId(String userId) {
    return _prefs.setString(USES_ID, userId);
  }
  Future<String> getUserId() async {
    return _prefs.getString(USES_ID) ?? '';
  }

  Future<bool> setDisplayName(String displayName) {
    return _prefs.setString(DISPLAY_NAME, displayName);
  }
  Future<String> getDisplayName() async {
    return _prefs.getString(DISPLAY_NAME) ?? '';
  }

  Future<bool> setQuickWord(String quickWord) {
    return _prefs.setString(QUICK_WORD, quickWord);
  }
  Future<String> getQuickWord() async {
    return _prefs.getString(QUICK_WORD) ?? '';
  }

  Future<bool> setAccountType(String accountType) {
    return _prefs.setString(ACCOUNT_TYPE, accountType);
  }
  Future<String> getAccountType() async {
    return _prefs.getString(ACCOUNT_TYPE) ?? '';
  }

  Future<bool> setMailAddress(String mailAddress) {
    return _prefs.setString(MAIL_ADDRESS, mailAddress);
  }
  Future<String> getMailAddress() async {
    return _prefs.getString(MAIL_ADDRESS) ?? '';
  }

  Future<bool> setSumStudyTime(Duration sumStudyTime) {
    return _prefs.setInt(SUM_STUDY_TIME, sumStudyTime.inMinutes);
  }
  Duration getSumStudyTime() {
    return new Duration(minutes: _prefs.getInt(SUM_STUDY_TIME) ?? 0);
  }

  Future<bool> setRegistrationDate(DateTime date) {
    return _prefs.setInt(REGISTRATION_DATE, date.millisecondsSinceEpoch);
  }
  DateTime getRegistrationDate() {
    return DateTime.fromMillisecondsSinceEpoch(_prefs.getInt(REGISTRATION_DATE) ?? 0);
  }

  // Future<void> setIdToken(String idToken) {
  //   return _secureStorage.write(key: ID_TOKEN, value: idToken);
  // }
  // Future<String> getIdToken() {
  //   return _secureStorage.read(key: ID_TOKEN) ?? '';
  // }
}