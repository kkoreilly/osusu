# Osusu

Osusu is an app for getting recommendations on what meals to eat in a group based on the ratings of each member of the group, and the cost, effort, healthiness, and recency of the meal.

To run it locally, you can follow the [Cogent Core installation instructions](https://www.cogentcore.org/core/setup/install), clone the repository, and then run the following commands:

```sh
rqlited -node-id=1 data/ # Run the database; see https://rqlite.io/docs/quick-start/
# In a new terminal tab/window:
cd cmd/osusu
core run
```

A web version will be deployed soon.
