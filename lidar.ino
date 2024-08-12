#include <Wire.h>
#include "SparkFun_VL53L1X.h" //Click here to get the library: http://librarymanager/All#SparkFun_VL53L1X

SFEVL53L1X distanceSensor;
int threshold = 10000; // Threshold to detect a peak

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

#define ARRAY_SIZE 20 // Define the size of the array to store signal rates

int signalRates[ARRAY_SIZE]; // Array to store signal rates
int currentIndex = 0;        // Index to keep track of current position in the array
int minSignalRate, maxSignalRate; // Variables to store minimum and maximum signal rates
bool minDetected = false;    // Flag to indicate if a minimum is detected
bool maxDetected = false;    // Flag to indicate if a maximum is detected

void loop(void)
{
  distanceSensor.startRanging();
  while (!distanceSensor.checkForDataReady())
  {
    delay(1);
  }

  int signalRate = distanceSensor.getSignalRate();
  //Serial.println(signalRate);

  // Store the signal rate in the array
  signalRates[currentIndex] = signalRate;
  currentIndex++;

  // If the array is full, process the data
  if (currentIndex >= ARRAY_SIZE)
  {
    minSignalRate = signalRates[0];
    maxSignalRate = signalRates[0];
    minDetected = false;
    maxDetected = false;

    // Find minimum and maximum values in the array
    for (int i = 1; i < ARRAY_SIZE; i++)
    {
      if (signalRates[i] < 0)
      {
        minSignalRate = signalRates[i];
        minDetected = true;
      }
      if (signalRates[i] > threshold)
      {
        maxSignalRate = signalRates[i];
        maxDetected = true;
      }
    }

    // Check if both a maximum and a minimum were detected in sequence
    bool specialCaseDetected = false;
    if (maxDetected && minDetected)
    {
      bool maxFirst = false;
      for (int i = 0; i < ARRAY_SIZE; i++)
      {
        if (signalRates[i] > threshold)
        {
          maxFirst = true;
        }
        if (maxFirst && signalRates[i] < 0)
        {
          specialCaseDetected = true;
          break;
        }
      }
    }

    // Handle the special case where maxDetected occurs before minDetected
    if (specialCaseDetected)
    {
    Serial.print("0");
    }
if (!specialCaseDetected)
    {
    Serial.print("1");
    }


    // Reset the index to overwrite the old data
    currentIndex = 0;
  }
}
