from locust import HttpUser, task, between
import uuid
import random
import string

# Utility helpers ------------------------------

def random_string(length=10):
    return ''.join(random.choices(string.ascii_lowercase, k=length))

def random_password():
    return ''.join(random.choices(string.ascii_letters + string.digits, k=12))

def random_fullname():
    first = random_string(6).capitalize()
    last = random_string(7).capitalize()
    return f"{first} {last}"

def random_sentence(words=10):
    return " ".join(random_string(random.randint(4,10)) for _ in range(words)) + "."

# Locust User Class -----------------------------

class BlogUser(HttpUser):
    wait_time = between(1, 3)

    def on_start(self):
        """
        Runs automatically when a simulated user starts.
        Creates:
           1) Random user
           2) Login
           3) Stores JWT token
        """
        # RANDOMIZED USER INFO
        self.username = random_string(12)
        self.password = random_password()
        self.name = random_fullname()

        # STEP 1 — REGISTER
        register_payload = {
            "name": self.name,
            "username": self.username,
            "password": self.password
        }

        with self.client.post(
            "/api/v1/users",
            json=register_payload,
            headers={"Content-Type": "application/json"},
            catch_response=True
        ) as resp:
            if resp.status_code not in [200, 201]:
                resp.failure("REGISTER FAILED: " + resp.text)
            else:
                resp.success()

        # STEP 2 — LOGIN
        login_payload = {
            "username": self.username,
            "password": self.password
        }

        with self.client.post(
            "/api/v1/auth/login",
            json=login_payload,
            headers={"Content-Type": "application/json"},
            catch_response=True
        ) as resp:
            if resp.status_code != 200:
                resp.failure("LOGIN FAILED: " + resp.text)
                return

            data = resp.json()
            self.token = data.get("token")

            if not self.token:
                resp.failure("LOGIN FAILED — No token returned")
            else:
                resp.success()

    @task
    def create_blog(self):
        """
        STEP 3 — Create blog with random content.
        """
        blog_payload = {
            "content": random_sentence(random.randint(10, 30))
        }

        headers = {
            "Content-Type": "application/json",
            "Authorization": f"Bearer {self.token}",
        }

        with self.client.post(
            "/api/v1/blogs",
            json=blog_payload,
            headers=headers,
            catch_response=True
        ) as resp:
            if resp.status_code not in [200, 201]:
                resp.failure("CREATE BLOG FAILED: " + resp.text)
            else:
                resp.success()
