const { initializeApp } = require('firebase/app');
const { getDatabase } = require('firebase/database');
const sqlite3 = require('sqlite3').verbose();
const db = new sqlite3.Database('hn.db');

// TODO: Replace the following with your app's Firebase project configuration
// See: https://firebase.google.com/docs/web/learn-more#config-object
const firebaseConfig = {
  // ...
  // The value of `databaseURL` depends on the location of the database
  databaseURL: "https://hacker-news.firebaseio.com/",
};

// Initialize Firebase
const app = initializeApp(firebaseConfig);


// Initialize Realtime Database and get a reference to the service
const database = getDatabase(app);
console.log(database);
