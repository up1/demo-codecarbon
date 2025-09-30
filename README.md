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
* Reduce dependencies
* Reduce startup time
* Reduce memory

With JAR
```
$./mvnw package
$ls -lh ./target/demo-java-0.0.1-SNAPSHOT.jar

$java -jar ./target/demo-java-0.0.1-SNAPSHOT.jar
```

With GraalVM (Don't need JRE)
```
$cd demo-java
$./mvnw -Pnative native:compile
$ls -lh ./target/demo-java
-rwxr-xr-x  1 somkiatpuisungnoen  staff    74M Sep 29 23:31 ./target/demo-java

$./target/demo-java
```