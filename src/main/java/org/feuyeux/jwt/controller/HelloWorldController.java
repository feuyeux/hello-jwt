package org.feuyeux.jwt.controller;

import org.reactivestreams.Publisher;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.CrossOrigin;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import reactor.core.publisher.Flux;

import java.time.Duration;
import java.util.List;
import java.util.Random;
import java.util.stream.IntStream;

import static java.util.stream.Collectors.toList;

@RestController
@CrossOrigin()
public class HelloWorldController {
    private static final Random random = new Random();

    @RequestMapping({"hello"})
    public String hello() {
        return "Hello World";
    }

    @GetMapping(value = "hello-stream", produces = MediaType.TEXT_EVENT_STREAM_VALUE)
    public Publisher<String> getHellos() {
        List<String> ids = IntStream.range(0, 3)
                .mapToObj(i -> getRandomId())
                .collect(toList());
        return Flux.fromIterable(ids).delayElements(Duration.ofMillis(50));
    }

    static String getRandomId() {
        int i = random.nextInt(5);
        return String.valueOf(i);
    }
}
