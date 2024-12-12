import newman from "newman";

newman.run({

    collection: "./auth_api_test.postman_collection.json",
    environment: "./auth_api_test.postman_environment.json",
    envVar: [
        {
            "key": "baseUrl",
            "value": process.env.BASE_URL
        },
        {
            "key": "secretTestClient",
            "value": process.env.SECRET_TEST_CLIENT
        },
        {
            "key": "randomKSUID",
            "value": process.env.RANDOM_KSUID
        }
    ],
    reporters: "cli"
},function (err) {
    if (err) {
        console.error(err?.message);
        throw err;
        return;
    }
    console.log('collection run complete!');
    process.exit(0);
});