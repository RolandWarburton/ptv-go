const fs = require('fs');
const crypto = require('crypto');

const devId = process.env.DEV_ID;
const key = process.env.KEY;
const baseURL = 'http://timetableapi.ptv.vic.gov.au';

function getUrl(request) {
  debugger;
  request = request + (request.includes('?') ? '&' : '?');
  const raw = request + `devid=${devId}`;
  const hmac = crypto.createHmac('sha1', key);
  hmac.update(raw);
  const signature = hmac.digest('hex');
  const url = `${baseURL}${raw}&signature=${signature}`;
  console.log(url);
  return url;
}

async function printFormattedDate(dateObj) {
  const formattedDate = dateObj.toLocaleString('en-GB', {
    // day: 'numeric',
    // month: 'numeric',
    // year: 'numeric',
    hour: 'numeric',
    minute: 'numeric',
    hour12: true
  });
  console.log(formattedDate);
}

async function getDepartures(stop_id, route_id, queryParams) {
  const requestString = `/v3/departures/route_type/0/stop/${stop_id}/route/${route_id}${queryParams}`;
  const url = getUrl(requestString);

  const res = await fetch(url);
  if (!res.ok) {
    throw new Error(`HTTP error! Status: ${res.status}`);
  }
  const data = await res.json();

  const departures = data.departures;
  fs.writeFileSync('output.json', JSON.stringify(departures));
  console.log(departures);
  return departures;
}

async function getNextDepartureTowards(departures, direction_id, count) {
  const now = new Date();

  // the departures to be returned
  const validDepartures = [];

  // track the number of results we want to return
  let i = 0;

  // iterate until you find a departure in the future
  for (const departure of departures) {
    const departureDate = new Date(departure.scheduled_departure_utc);
    if (!departureDate) {
      throw new Error('failed');
    }
    // if the train already departed ignore it
    if (departureDate.getTime() < now.getTime()) {
      continue;
    }

    // if the train is not going towards the nominated direction
    if (departure.direction_id !== direction_id) {
      continue;
    }

    validDepartures.push({ ...departure, index: i });
    i++;
    if (i == count) {
      return validDepartures;
    }
  }
  return validDepartures;
}

async function main() {
  // get the Bayswater departures
  //route_id = 2         belgrave line
  // stop_id = 1016      bayswater
  const departures = await getDepartures(1016, 2, '?expand=All&include_geopath=true');

  // direction_id = 1    going towards city
  const nextBWDepartures = await getNextDepartureTowards(departures, 1, 2);
  printFormattedDate(new Date(nextBWDepartures[0].scheduled_departure_utc));
  printFormattedDate(new Date(nextBWDepartures[1].scheduled_departure_utc));
  console.log(nextBWDepartures);
}

main();
