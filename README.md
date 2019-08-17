# Twitter kafka connector

This is a small go application that allows you to stream tweets into a kafka cluster.

## Configuration

All settings can be set via a config.yml in /etc/kafka_twitter_connector, ./ or enviroment variables. Note that env vars have
to be upper case.

### Kafka settings
- `kafka_brokers`: The kafka brokers, to which the tweets should be streamed. By default "localhost:9092"
- `schema_registry`: The schema registry host. By default "localhost:8081
- `tweet_topic`: The topic to which the should be streamed. By default "twees"

### Twitter settings
You will need to request a developer account and create an app to get the following oauth settings.

- `client_key`: The client API key.
- `client_secret`: The client Secret key.
- `access_token`: The access token key.
- `access_secret`: The access secret key.

### Twitter filters
You use the following settings to filter the stream. Note that all filters are concatinated by OR, this is a limitiation of
the public streaming API. Example tracking "dog" and "cat" means all posts with containg "dog" or "cat" will be streamed to
kafka.

You also need to set at least one filter in order for the stream to work.

If you use enviroment variables, use a space to seperate multiple items, e.g. `export keyword="kitten cat"

- `keywords`: Track the given keywords.
