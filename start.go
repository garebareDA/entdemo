package main

import (
	"context"
	"entdemo/ent"
	"entdemo/ent/user"
	"entdemo/ent/car"
	"entdemo/ent/group"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatal("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
}

func CreateUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.
		Create().
		SetAge(30).
		SetName("a8m").
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %v", err)
	}
	log.Println("user was created", u)
	return u, nil
}

func QueryUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.Query().Where(user.Name("a8m")).Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying user: %v", err)
	}
	log.Println("user returned: ", u)
	return u, nil
}

func CreateCars(ctx context.Context, client *ent.Client) (*ent.User, error) {
	tesla, err := client.Car.Create().SetModel("Tesla").SetRegisteredAt(time.Now()).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating car: %v", err)
	}
	log.Println("user returned: ", tesla)

	ford, err := client.Car.Create().SetModel("Ford").SetRegisteredAt(time.Now()).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating car: %v", err)
	}
	log.Println("car was created: ", ford)

	a8m, err := client.User.Create().SetAge(30).SetName("a8m").AddCars(tesla, ford).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating car: %v", err)
	}
	log.Println("user was created: ", a8m)

	return a8m, nil
}

func QueryCars(ctx context.Context, a8m *ent.User) error {
	cars, err := a8m.QueryCars().All(ctx)
	if err != nil {
		return fmt.Errorf("error querying cars: %v", err)
	}
	log.Println("cars: ", cars)

	ford, err := a8m.QueryCars().Where(car.Model("Ford")).Only(ctx)
	if err != nil {
		return fmt.Errorf("error querying cars: %v", err)
	}
	log.Println("cars: ",ford)
	return nil
}

func QueryCarUsers(ctx context.Context, a8m *ent.User) error {
	cars, err := a8m.QueryCars().All(ctx)
	if err != nil {
		return fmt.Errorf("error querying cars: %v", err)
	}

	for _, c := range cars {
		owner, err := c.QueryOwner().Only(ctx)
		if err != nil {
			return fmt.Errorf("failed querying car %q owner: %w", c.Model, err)
		}
		log.Printf("car %q owner: %q\n", c.Model, owner)
	}

	return nil
}

func CreateGraph(ctx context.Context, client *ent.Client) error {
	a8m, err := client.User.Create().SetAge(30).SetName("Ariel").Save(ctx)
	if err != nil {
		return err
	}
	neta, err := client.User.Create().SetAge(28).SetName("Neta").Save(ctx)
	if err != nil {
		return err
	}

	err = client.Car.Create().SetModel("Tesla").SetRegisteredAt(time.Now()).SetOwner(a8m).Exec(ctx)
	if err != nil {
		return err
	}

	err = client.Car.Create().SetModel("Mazda").SetRegisteredAt(time.Now()).SetOwner(a8m).Exec(ctx)
	if err != nil {
		return err
	}

	err = client.Car.Create().SetModel("Ford").SetRegisteredAt(time.Now()).SetOwner(a8m).Exec(ctx);
	if err != nil {
		return err
	}

	err = client.Group.Create().SetName("GitLab").AddUsers(neta, a8m).Exec(ctx)
	if err != nil {
		return err
	}

	err = client.Group.Create().SetName("GitHub").AddUsers(a8m).Exec(ctx)
	if err != nil {
		return err
	}

	log.Println("The graph was created successfully")
	return nil
}

func QueryGitHub(ctx context.Context, client *ent.Client) error {
	cars, err := client.Group.Query().Where(group.Name("GitHub")).QueryUsers().QueryCars().All(ctx)
	if err != nil {
		return fmt.Errorf("failed getting cars: %w", err)
	}
	log.Println("cars returned:", cars)

	return nil
}

func QueryArielCars(ctx context.Context, client *ent.Client) error {
	a8m := client.User.Query().Where(user.HasCars(), user.Name("Ariel")).OnlyX(ctx)
	cars, err := a8m.QueryGroups().QueryUsers().QueryCars().Where(car.Not(car.Model("Mazda"))).All(ctx)
	if err != nil {
		return fmt.Errorf("failed getting cars: %w", err)
	}
	log.Println("cars returned:", cars)
	return nil
}

func QueryGroupWithUsers(ctx context.Context, client *ent.Client) error {
	groups, err := client.Group.Query().Where(group.HasUsers()).All(ctx)
	if err != nil {
		return fmt.Errorf("failed getting groups: %w", err)
	}
	log.Println("groups returned:", groups)
	return nil
}