# Demo with Green Coding

## 1. Working with Python and codecarbon
* https://github.com/mlco2/codecarbon
```
$pip install codecarbon
$export CODECARBON_API_KEY=<your-api-key>
$codecarbon monitor

$cd demo-python
# List
$python demo01.py

# Set
$python demo01.py
```

## 2. Data layout of format
* From CSV to [Apache Parquet](https://parquet.apache.org/)

```
$pip install -r requirements.txt
$python demo-data-format.py
```

## 3. Serialization & APIs
* From JSON over REST to gRPC + Protobuf for hot paths (smaller payloads, faster)


## Java and GraalVM with Spring Boot
* https://www.graalvm.org/latest/reference-manual/native-image/guides/build-spring-boot-app-into-native-executable/
```
$cd demo-java
$./mvnw -Pnative native:compile
```