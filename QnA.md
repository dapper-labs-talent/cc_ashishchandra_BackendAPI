How long did this assignment take?
> It took about 3 days to complete as I had some clarifications I needed around the deployment environment, and setting up of tests.

What was the hardest part?
> Out of everything, the hardest part was ensuring it will run from anywhere. That is why I had to introduce docker into the mix and use dockerized postgres.

Did you learn anything new?
> Yes of course! I am fairly new to Golang as we only switched to using Golang (from C++ and Node.js) last year in my company. So every new problem that I solve teaches me something. From this assignment, I learnt some finer points of setting up unit tests using the Golang testing framework.

Is there anything you would have liked to implement but didn't have the time to?
> I think I completed everything that was assigned. It would have been nice to run some stress tests to determine how well the server holds up against a heavy load of users.

What are the security holes (if any) in your system? If there are any, how would you fix them?
> I doubt it. I deliberately went the route of using server based signing and verification keys that are based on ECDSA ES512 (SHA 512 hash) so unless someone gains access to the server where the application is running and gets hold of the private key that signs the JWT, it will be very difficulty to spoof a login.

Do you feel that your skills were well tested?
> Yes I think so. Like I said, I only switched to Golang last year (and love it!) so I like being given a challenge, and then solving it.