scheduler to automatically create promos for users who are having birthdays and will send them to the user's WhatsApp number, this system runs on **Docker** and **Kafka as a message broker** and **SMS gateway** from **zenziva.net**, you can use your own account or use an account that has been created in the **.env** file, with limited messages of course.

**Before running the system, change the phone and date_birth data according to what you want in the insert user data query in the init.sql file, on this query:**
_INSERT IGNORE INTO user (name, email, phone, date_birth)
VALUES ("Agung","agung@gmail.com","08978787687","1996-11-09");_

_note: the scheduler is triggered every minute and when a promo has been created today, the promo will not be created again_

## Run System

```bash
$ docker-compose up -d
```

## Flowchart

![alt text](https://github.com/dedihartono801/promo-scheduler/blob/master/scheduler_flowchart.png)
