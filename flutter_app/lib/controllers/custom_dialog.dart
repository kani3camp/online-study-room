import 'package:flutter/material.dart';

class CustomDialog {
  static void showLoadingDialog(BuildContext context, {String title=''}) {
    showDialog(
      barrierDismissible: false,
      context: context,
      builder: (_) {
        return AlertDialog(
          title: Text(title),
          content: SizedBox(
            height: 100,
            child: Center(
              child: SizedBox(
                width: 50,
                height: 50,
                child: CircularProgressIndicator(),
              ),
            ),
          ),
        );
      },
    );
  }

  static Future<void> showAlertDialog(
      BuildContext context,
      String message,
      {
        Function onOkPressed,
        Function onCancelPressed,
      })
  async {
    await showDialog(
        context: context,
        builder: (context) {
          return AlertDialog(
            title: Text(message),
            content: null,
            actions: <Widget>[
              // ボタン領域
              onCancelPressed != null ? FlatButton(
                child: Text("キャンセル"),
                onPressed: onCancelPressed,
              )
              : null,
              FlatButton(
                child: Text("OK"),
                onPressed: onOkPressed,
              ),
            ],
          );
        }
    );
  }
}
