import 'package:flutter/material.dart';

class LoadingDialog {
  static void show(BuildContext context, {String title=''}) {
    showDialog(
      context: context,
      builder: (_) {
        return AlertDialog(
          title: Text(title),
          content: FittedBox(
            child: CircularProgressIndicator(),
          ),
        );
      },
    );
  }
}