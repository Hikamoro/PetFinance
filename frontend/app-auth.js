// Простая логика регистрации / авторизации + локальная часть бюджета

let incomes = [];
let expenses = [];

// Auth
const authSection = document.getElementById('auth');
const appSection = document.getElementById('app');
const authMessage = document.getElementById('authMessage');

const regLogin = document.getElementById('regLogin');
const regPassword = document.getElementById('regPassword');
const loginLogin = document.getElementById('loginLogin');
const loginPassword = document.getElementById('loginPassword');

const authApiBase = '' // Оставить пустым, чтобы использовать тот же origin (http://localhost:8080)

// Ссылки на элементы DOM
const incomeAmount = document.getElementById('incomeAmount');
const incomeSource = document.getElementById('incomeSource');
const incomeComment = document.getElementById('incomeComment');

const expenseAmount = document.getElementById('expenseAmount');
const expenseCategory = document.getElementById('expenseCategory');
const expenseComment = document.getElementById('expenseComment');

const balanceEl = document.getElementById('balance');
const incomeTotal = document.getElementById('incomeTotal');
const expenseTotal = document.getElementById('expenseTotal');
const historyList = document.getElementById('history');

function setAuthMessage(msg, isError = true) {
  if (!authMessage) return;
  authMessage.textContent = msg;
  authMessage.style.color = isError ? '#b00' : '#080';
}

function showLogin() {
  document.getElementById('loginCard').style.display = 'block';
  document.querySelector('.auth-card').style.display = 'none';
  setAuthMessage('');
}

function showRegister() {
  document.getElementById('loginCard').style.display = 'none';
  document.querySelector('.auth-card').style.display = 'block';
  setAuthMessage('');
}

async function register() {
  setAuthMessage('');
  const login = regLogin.value.trim();
  const password = regPassword.value;

  if (!login || !password) {
    setAuthMessage('Введите логин и пароль');
    return;
  }

  const res = await fetch(`${authApiBase}/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ login, password })
  });

  if (!res.ok) {
    const text = await res.text();
    setAuthMessage(text || 'Ошибка регистрации');
    return;
  }

  setAuthMessage('Регистрация прошла успешно. Войдите.', false);
  showLogin();
}

async function login() {
  setAuthMessage('');
  const login = loginLogin.value.trim();
  const password = loginPassword.value;

  if (!login || !password) {
    setAuthMessage('Введите логин и пароль');
    return;
  }

  const res = await fetch(`${authApiBase}/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ login, password })
  });

  if (!res.ok) {
    const text = await res.text();
    setAuthMessage(text || 'Ошибка входа');
    return;
  }

  const data = await res.json();
  if (!data.api_hash) {
    setAuthMessage('Не удалось получить api_hash');
    return;
  }

  localStorage.setItem('api_hash', data.api_hash);
  localStorage.setItem('username', login);
  setAuthMessage('');
  showApp();
}

function showApp() {
  if (authSection) authSection.style.display = 'none';
  if (appSection) appSection.style.display = 'block';
  render();
}

function init() {
  const apiHash = localStorage.getItem('api_hash');
  if (apiHash) {
    showApp();
  } else {
    if (authSection) authSection.style.display = 'block';
    if (appSection) appSection.style.display = 'none';
  }
}

/* ДОХОД */

function addIncome() {
  const amount = Number(incomeAmount.value);
  const source = incomeSource.value;
  const comment = incomeComment.value;

  if (!amount) {
    alert('Введите сумму дохода');
    return;
  }

  incomes.push({
    id: Date.now(),
    amount,
    source,
    comment
  });

  clearIncome();
  render();
}

function clearIncome() {
  incomeAmount.value = '';
  incomeSource.value = '';
  incomeComment.value = '';
}

/* РАСХОД */

function addExpense() {
  const amount = Number(expenseAmount.value);
  const category = expenseCategory.value;
  const comment = expenseComment.value;

  if (!amount) {
    alert('Введите сумму расхода');
    return;
  }

  expenses.push({
    id: Date.now(),
    amount,
    category,
    comment
  });

  clearExpense();
  render();
}

function clearExpense() {
  expenseAmount.value = '';
  expenseCategory.value = '';
  expenseComment.value = '';
}

/* РЕНДЕР */

function render() {
  const incomeSum = incomes.reduce((s, i) => s + i.amount, 0);
  const expenseSum = expenses.reduce((s, e) => s + e.amount, 0);

  const balance = incomeSum - expenseSum;

  balanceEl.innerText = balance.toLocaleString() + ' ₽';
  incomeTotal.innerText = incomeSum.toLocaleString() + ' ₽';
  expenseTotal.innerText = expenseSum.toLocaleString() + ' ₽';

  renderHistory();
}

function renderHistory() {
  historyList.innerHTML = '';

  const all = [
    ...incomes.map(i => ({ ...i, type: 'income' })),
    ...expenses.map(e => ({ ...e, type: 'expense' }))
  ];

  all.sort((a, b) => b.id - a.id);

  all.forEach(item => {
    const li = document.createElement('li');

    if (item.type === 'income') {
      li.innerHTML = `
        <span>➕ ${item.source || ''} ${item.comment ? '— ' + item.comment : ''}</span>
        <strong style="color:green">+${item.amount} ₽</strong>
      `;
    } else {
      li.innerHTML = `
        <span>➖ ${item.category || ''} ${item.comment ? '— ' + item.comment : ''}</span>
        <strong style="color:red">-${item.amount} ₽</strong>
      `;
    }

    historyList.appendChild(li);
  });
}

init();
