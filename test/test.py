# should be launched with locust -f .\test.py
import random
import string
from locust import HttpUser, TaskSet, task, between

def generate_random_username(length=8):
    """Генерирует случайное имя пользователя."""
    return ''.join(random.choices(string.ascii_lowercase + string.digits, k=length))

class UserBehavior(TaskSet):
    users = []  # Общий список пользователей для всех экземпляров

    def on_start(self):
        self.username = generate_random_username()
        while self.username in self.users:
            self.username = generate_random_username()
        response = self.client.post("/api/auth", json={"username": self.username, "password": "password"})
        if response.status_code == 200:
            try:
                result = response.json()
                self.token = result["token"]
                self.users.append(self.username)
            except ValueError:
                print(f"Ошибка при парсинге JSON: {response.text}")
                self.token = None

        else:
            print(f"Ошибка аутентификации: {response.status_code}, {response.text}")
            self.token = None

    @task(3)
    def get_info(self):
        if self.token != None:
            self.client.get("/api/info", headers={"Authorization": f"Bearer {self.token}"})

    @task(1)
    def send_coin(self):
        if self.token != None and self.users:
            to_user = random.choice(self.users)
            self.client.post("/api/sendCoin", json={"toUser": to_user, "amount": 10}, headers={"Authorization": f"Bearer {self.token}"})

    @task(2)
    def buy_item(self):
        if self.token != None:
            self.client.get("/api/buy/pen", headers={"Authorization": f"Bearer {self.token}"})

class WebsiteUser(HttpUser):
    tasks = [UserBehavior]
    wait_time = between(1, 1)
