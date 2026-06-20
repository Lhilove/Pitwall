import fastf1
import requests
import time

# enable caching so FastF1 doesn't re-download data every run
fastf1.Cache.enable_cache('cache')

# load a real session — change year/round/session as needed
session = fastf1.get_session(2024, 'Bahrain', 'R')
session.load()

API_URL = "http://localhost:8080/telemetry"

# loop through each driver in the session
for driver in session.drivers:
    laps = session.laps.pick_driver(driver)

    for _, lap in laps.iterlaps():
        car_data = lap.get_car_data()

        # send a sample of points per lap, not every single one (would be thousands)
        for i in range(0, len(car_data), 20):
            point = car_data.iloc[i]

            payload = {
                "driver": str(driver),
                "lap": int(lap['LapNumber']),
                "speed": int(point['Speed']),
                "throttle": int(point['Throttle']),
                "brake": int(point['Brake']),
                "gear": int(point['nGear']),
            }

            try:
                requests.post(API_URL, json=payload, timeout=2)
            except requests.exceptions.RequestException as e:
                print("error sending:", e)

            time.sleep(0.05)  # small delay so it streams instead of dumping instantly

print("done feeding telemetry")