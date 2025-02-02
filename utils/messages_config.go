package utils

// MessagesConfig ...
var MessagesConfig = `{
  "00001I": {
    "statusCode": 200,
    "messageCode": "00001I",
    "message": {
      "ja": "すでにお気に入り登録済みです。",
      "en": "Already liked."
    }
  },
  "00001E": {
    "statusCode": 400,
    "messageCode": "00001E",
    "message": {
      "ja": "ヘッダーチェック処理にてエラーが発生しました。",
      "en": "header parameter is invalid."
    }
  },
  "00002E": {
    "statusCode": 400,
    "messageCode": "00002E",
    "message": {
      "ja": "クライアントIDチェック処理にてエラーが発生しました。",
      "en": "ClientID parameter is invalid."
    }
  },
  "10001E": {
    "statusCode": 500,
    "messageCode": "10001E",
    "message": {
      "ja": "DBへのデータ取得時にエラーが発生しました。",
      "en": "DB select error."
    }
  },
  "10002E": {
    "statusCode": 500,
    "messageCode": "10002E",
    "message": {
      "ja": "オブジェクトの変換に失敗しました。",
      "en": "object mapping error."
    }
  },
  "10003E": {
    "statusCode": 500,
    "messageCode": "10003E",
    "message": {
      "ja": "DBへのデータ保存時にエラーが発生しました。",
      "en": "DB update error."
    }
  }
}`
