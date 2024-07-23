const { initializeApp } = require("firebase/app");
const { getDatabase, ref, child, get } = require("firebase/database");
const sqlite3 = require("sqlite3").verbose();
const db = new sqlite3.Database("hn.db");

(async () => {
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

  const dbRef = ref(getDatabase());
  try {
    const snapshot = await get(child(dbRef, "/v0/maxitem"));
    if (snapshot.exists()) {
      console.log(snapshot.val());
    } else {
      console.log("No data available");
    }
  } catch (error) {
    console.error(error);
  }
})();
