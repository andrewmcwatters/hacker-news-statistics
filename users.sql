SELECT `by`
FROM (
  SELECT
    `by`,
    MIN(timestamp) AS timestamp
  FROM `bigquery-public-data.hacker_news.full`
  WHERE `by` IS NOT NULL
  GROUP BY `by`
)
ORDER BY timestamp;
