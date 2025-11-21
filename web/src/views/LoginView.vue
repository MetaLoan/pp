<template>
  <div class="login-container">
    <div class="login-box">
      <h2>Login to PP</h2>
      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <label>Email</label>
          <input type="email" v-model="email" required />
        </div>
        <div class="form-group">
          <label>Password</label>
          <input type="password" v-model="password" required />
        </div>
        <button type="submit" :disabled="loading">
          {{ loading ? 'Logging in...' : 'Login' }}
        </button>
      </form>
      <p class="register-link">
        Don't have an account? <a href="#" @click.prevent="handleRegister">Register (Demo)</a>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';

const router = useRouter();
const authStore = useAuthStore();

const email = ref('');
const password = ref('');
const loading = ref(false);

const handleLogin = async () => {
  loading.value = true;
  try {
    await authStore.login(email.value, password.value);
    router.push('/trade');
  } catch (error) {
    alert('Login failed: ' + (error.response?.data?.error || error.message));
  } finally {
    loading.value = false;
  }
};

const handleRegister = async () => {
  // Quick demo registration
  const randomId = Math.floor(Math.random() * 1000);
  const demoEmail = `user${randomId}@example.com`;
  const demoPass = 'password123';
  
  try {
    await authStore.register(demoEmail, demoPass, `User${randomId}`);
    alert(`Registered demo account:\nEmail: ${demoEmail}\nPassword: ${demoPass}`);
    email.value = demoEmail;
    password.value = demoPass;
  } catch (error) {
    alert('Registration failed');
  }
};
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #131722;
  color: white;
}

.login-box {
  background: #1e222d;
  padding: 40px;
  border-radius: 8px;
  width: 100%;
  max-width: 400px;
  box-shadow: 0 4px 6px rgba(0,0,0,0.3);
}

h2 {
  text-align: center;
  margin-bottom: 30px;
}

.form-group {
  margin-bottom: 20px;
}

label {
  display: block;
  margin-bottom: 8px;
}

input {
  width: 100%;
  padding: 12px;
  background: #2b2b43;
  border: 1px solid #43435c;
  color: white;
  border-radius: 4px;
  box-sizing: border-box;
}

button {
  width: 100%;
  padding: 12px;
  background-color: #2962FF;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
}

button:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.register-link {
  text-align: center;
  margin-top: 20px;
  font-size: 14px;
}

a {
  color: #2962FF;
  text-decoration: none;
}
</style>
