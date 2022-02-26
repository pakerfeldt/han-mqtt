#!/usr/bin/env node

const mqtt = require('mqtt')
const config = require('./config.js').parse()
const { SerialPort } = require('serialport')
const { ReadlineParser } = require('@serialport/parser-readline')
const obis = require('./obismap').description

const regexp = /([0-9:.-]+)\(([\d.]+)\*(\w+)\)/

const client = mqtt.connect(config.mqtt.url, config.mqtt.options)
const port = new SerialPort(config.serial)
const parser = port.pipe(new ReadlineParser({ delimiter: '\r\n' }))

let topicPrefix = (config.mqtt.topicPrefix !== undefined && config.mqtt.topicPrefix !== "") ? config.mqtt.topicPrefix + "/" : ""

parser.on('data', function (line) {
    let match = line.match(regexp)
    if (match !== null) {
        let message = line
        if (!config.sendUnparsed) {
            message = JSON.stringify({
                id: match[1],
                value: parseFloat(match[2]),
                unit: match[3],
                description: obis[match[1]]
            })
        }
        client.publish(topicPrefix + match[1], message)
    }

})
