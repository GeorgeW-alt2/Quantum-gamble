#include <Wire.h>
#include "SparkFun_VL53L1X.h" //Click here to get the library: http://librarymanager/All#SparkFun_VL53L1X

SFEVL53L1X distanceSensor;
void setup(void)
{
  pinMode(LED_BUILTIN, OUTPUT);
  digitalWrite(LED_BUILTIN, HIGH);
  Wire.begin();
  Serial.begin(115200);
  if (distanceSensor.begin() != 0) // Begin returns 0 on a good init
  {
    while (1)
      ;
  }
}

void loop(void)
{
  distanceSensor.startRanging();
  while (!distanceSensor.checkForDataReady())
  {
    delay(1);
  }

  int signalRate = distanceSensor.getSignalRate();
  //Serial.println(abs(signalRate));
  if (abs(signalRate) > 22000){
    Serial.print("1");
    }
  if (abs(signalRate) > 15000 && abs(signalRate) < 22000 ){
    Serial.print("0");
    }
  delay(10);
}
