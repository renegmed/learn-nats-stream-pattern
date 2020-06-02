## Explore and Learn NATS Streaming Integration Pattern ##


Default configuration. One publisher and two subscribers receiving messages at the same time.

To run:

1. Create and deploy containers

    $ make up

2. On new terminal, monitor sub1. (see Makefile for others)

    $ make tail-sub1

3. Post a transaction

    $ make post

