# rest-git

REST interface around your git repo

## Sample config file

```toml
RefreshInterval = "5m"

[Repositories]

[Repositories.Firebase]
name = "firebase"
url = "https://github.com/firebase/firebase-js-sdk.git"

[Repositories.Symfony]
name = "symfony"
url = "https://github.com/symfony/symfony.git"

[Repositories.Logstash]
name = "Logstash"
url = "https://github.com/elastic/logstash.git"

[Repositories.MechanicalSoup]
name = "MechanicalSoup"
url = "https://github.com/MechanicalSoup/MechanicalSoup.git"

[Repositories.the-book]
name = "the-book"
url = "https://github.com/trimstray/the-book-of-secret-knowledge.git"
```

## Available endpoints

### /list

Lists all tracked repositories

### /fetch

Update all checkouts with the recent changes

### /{repo-id}/fetch

Update selected repository with the recent changes

### /{repo-id}/branches

### /{repo-id}/tags

Lists tags / branches for the particular repo. Sample output:

```json
{
  "branches": [
    {
      "Name": "v2.0.0",
      "Sha": "867975dc84167c484fbc9e300061ca9d267fade4"
    },
    {
      "Name": "v2.0.0-RC1",
      "Sha": "6d8932f53849ab1d073490135b847170afd46b76"
    },
    {
      "Name": "v2.0.0-RC2",
      "Sha": "204f143d21bb9f83db2a956b7ed00f46e7b6487b"
    }
  ]
}
```

## TODO

- [ ] add git log
- [ ] add git diff
- [ ] support in memory checkout
- [ ] start webserver before checkouts are ready
- [ ] fix reload config endpoint