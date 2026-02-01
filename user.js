const fs = require("fs");
const path = require("path");
const sqlite3 = require("sqlite3").verbose();

const db = new sqlite3.Database("hacker-news.db");
const csvPath = path.join(__dirname, "bq-results-20260201-205857-1769979839843.csv");

// Create user table if it doesn't exist
const createTableSQL = `CREATE TABLE IF NOT EXISTS user (
  id TEXT PRIMARY KEY,
  created INTEGER,
  karma INTEGER,
  updated_at TEXT
)`;
db.run(createTableSQL);

const fetch = (...args) => import('node-fetch').then(({default: fetch}) => fetch(...args));

(async () => {
  const lines = fs.readFileSync(csvPath, "utf-8").split("\n").filter(Boolean);
  const usernames = lines.slice(1); // skip header
  for (const username of usernames) {
    try {
      const res = await fetch(`https://hacker-news.firebaseio.com/v0/user/${username}.json`);
      if (!res.ok) continue;
      const user = await res.json();
      if (!user) continue;
      const { id, created, karma } = user;
      const updated_at = new Date().toISOString();
      db.run(
        `INSERT INTO user (id, created, karma, updated_at) VALUES (?, ?, ?, ?)
         ON CONFLICT(id) DO UPDATE SET created=excluded.created, karma=excluded.karma, updated_at=excluded.updated_at`,
        [id, created, karma, updated_at]
      );
      console.log(`Upserted user: ${id}`);
    } catch (e) {
      console.error(`Failed for user ${username}:`, e);
    }
  }
  db.close();
})();
