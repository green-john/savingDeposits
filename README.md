# Saving Deposits

Manage saving deposits

## Install

Go 1.11+

## Usage

The application is written with a Go backend and a Javascript (vue.js) frontend.
Make sure `baseUrl` is set correctly in `frontend/src/components/http.js`, otherwise
requests to the server will fail.

The server takes the database host, name and user from env variables, specifically:

```
SAVINGS_DB_HOST
SAVINGS_DB_NAME
SAVINGS_DB_USER
```

Postgresql is used as a database. Make sure you `createdb` before starting the app.

See `scripts/run.sh` for an example on how to start the server. `scripts/rentals-cli`
is a compiled binary that can be used directly to run the server. Otherwise, you can install
go > 1.11 (New modules are used) and build all the project. To obtain a binary, cd into
`cmd/rentals-cli` and run `go build ./...`.

Scripts in `scripts` DO NOT WORK if you are outside the directory. Make sure you cd into
it.

To run all tests, run `scripts/all_tests.sh`. As previously stated, make sure you
are inside the `scripts` directory.

The frontend is run using `vue-cli`. `cd` into the `frontend` folder, install all dependencies
with `npm install --saveDev` and then run `npm run serve`. This will bring up a dev server
for the frontend.


## Creating a first user

In order to start interacting with the system a first user must be created. The api only allows
to create clients, so in order to get admin access, you must create a client user, and then
change the role directly in the database. This way you will have an admin that can then
be used to create realtors and more admins.

## Docs

The api is documented using [Open API 2.0](https://swagger.io/specification/). See `docs/api.yml`.
To view it in the browser install `redoc-cli` (`npm install -g redoc-cli`) and then run
`redoc-cli serve docs/api.yml`.

## Spec

- User must be able to create an account and log in. (If a mobile application, this means that more users can use
the app from the same phone).
- When logged in, a user can see, edit and delete saving deposits he entered.
- Implement at least three roles with different permission levels: a regular user would only be able to CRUD on their
owned records, a user manager would be able to CRUD users, and an admin would be able to CRUD all records and users.
- A saving deposit is identified by a bank name, account number, an initial amount saved (currency in USD), start date, end date, interest percentage per year and taxes percentage.
- The interest could be positive or negative. The taxes are only applied over profit.
- User can filter saving deposits by the amount (minimum and maximum), bank name and date.
- User can generate a revenue report for a given period, that will show the gains and losses
from interest and taxes for each deposit. The amount should be green or red if respectively
it represents a gain or loss. The report should show the sum of profits and losses at the bottom
for that period. 
- The computation of profit/loss is done on a daily basis. Consider that a year is 360 days. 
- REST API. Make it possible to perform all user actions via the API, including authentication (If a mobile application and you don’t know how to create your own backend you can use Firebase.com or similar services to create the API).
- In any case, you should be able to explain how a REST API works and demonstrate that by creating functional tests that use the REST Layer directly. Please be prepared to use REST clients like Postman, cURL, etc. for this purpose.
- If it’s a web application, it must be a single-page application. All actions need to be done client side using AJAX, refreshing the page is not acceptable. (If a mobile application, disregard this).
- Functional UI/UX design is needed. You are not required to create a unique design, however, do follow best practices to make the project as functional as possible.
- Bonus: unit and e2e tests.

## TODO

- Add logging to db errors (Pretty much done)
- Find a better structure for the files. Right now there are two
places where routes are added. (Almost done)
- Add ctrl+z cancellation signal.
- Add user frontend.
- Add some vue material library so shit looks better.
- Figure out how to pass values dynamically.

## Next actions:

- ~Frontend to CRUD deposits.~
- ~Frontend to search deposits~
- Frontend for report.
- Backend search.
- Backend report.
- Fix rates to be any number not just decimal
- Vselect for users.
- Fix permissions for getting deposits for a user.

