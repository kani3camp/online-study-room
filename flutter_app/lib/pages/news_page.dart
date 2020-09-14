import 'dart:convert';

import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

class NewsPage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    Future<List> futureList = fetchNewsList();

    return Scaffold(
      appBar: AppBar(
        title: Center(
          child: Text('お知らせ')
        ),
      ),
      body: FutureBuilder<List>(
        future: futureList,
        builder: (context, snapshot) {
          if (snapshot.hasData) {
            return ListView.separated(
              padding: const EdgeInsets.all(8),
              itemCount: snapshot.data.length,
              itemBuilder: (BuildContext context, int index) {
                return Container(
                  child: ListTile(
                    title: Text(snapshot.data[index].newsBody.title),
                    subtitle: Text(snapshot.data[index].newsBody.textBody),
                  ),
                );
              },
              separatorBuilder: (BuildContext context, int index) => const Divider(),
            );
          } else if (snapshot.hasError) {
            print(snapshot.error);
          } else {
          }
          return Center(child: CircularProgressIndicator());
        },
      ),
    );
  }
}

Future<List> fetchNewsList() async {
  Map<String, String> queryParams = {
    'num_news': '10'
  };
  Uri uri = Uri.https('us-central1-online-study-room-f1f30.cloudfunctions.net', '/News', queryParams);
  final response = await http.get(uri);
  if (response.statusCode == 200) {
    NewsResponse roomsResponse = NewsResponse.fromJson(json.decode(utf8.decode(response.bodyBytes)));
    if (roomsResponse.result == 'ok') {
      return roomsResponse.newsList;
    } else {
      throw Exception('Failed to load news list: ' + roomsResponse.message);
    }
  } else {
    throw Exception('http request failed');
  }
}

class NewsResponse {
  final String result;
  final String message;
  final List<News> newsList;

  NewsResponse({this.result, this.message, this.newsList});

  factory NewsResponse.fromJson(Map<String, dynamic> json) {
    return NewsResponse(
      result: json['result'] as String,
      message: json['message'] as String,
      newsList: (json['news_list'] as List<dynamic>).map((i) => News.fromJson(i)).toList()
    );
  }
}

class News {
  final String newsId;
  final NewsBody newsBody;

  News({this.newsId, this.newsBody});

  factory News.fromJson(Map<String, dynamic> json) {
    return News(
      newsId: json['news_id'] as String,
      newsBody: NewsBody.fromJson(json['news_body'])
    );
  }
}

class NewsBody {
  final DateTime created;
  final DateTime updated;
  final String title;
  final String textBody;

  NewsBody({this.created, this.updated, this.title, this.textBody});

  factory NewsBody.fromJson(Map<String, dynamic> json) {
    return NewsBody(
      created: DateTime.parse(json['created']),
      updated: DateTime.parse(json['updated']),
      title: json['title'] as String,
      textBody: json['text_body'] as String
    );
  }
}