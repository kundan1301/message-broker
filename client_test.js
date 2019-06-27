const mqtt = require('mqtt')
const client  = mqtt.connect('mqtt://localhost:1883',{username:'kundan',
    password:'kundan',clientId:'my_name',clean:false, qos:2}
)

client.on('connect', function () {
    console.log('connected')    
})

client.on('error',function(error){
    console.log('error',error);
})