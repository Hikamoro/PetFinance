const apiKey = localStorage.getItem('api_hash');
  if (!apiKey) {
    window.location.href = '/auth';
  }

  document.addEventListener('DOMContentLoaded', () => {
    const username = localStorage.getItem('username') || 'User';
    const logo = document.querySelector('.logo');
    if (logo) logo.textContent = username;

    loadDashboard();
  });

  function logout() {
    localStorage.removeItem('api_hash');
    localStorage.removeItem('username');
    window.location.href = '/auth';
  }

  function formatMoney(value) {
    return '$' + Number(value).toLocaleString('en-US');
  }

  async function postJson(url, data) {
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      const text = await response.text();
      throw new Error(text || `${response.status} ${response.statusText}`);
    }

    return response.json().catch(() => null);
  }

  async function fetchBalance() {
    const result = await postJson('/getBalance', { api_hash: apiKey });
    return Number(result?.balance || 0);
  }

  async function fetchIncomes() {
    const result = await postJson('/checkIncome', { api_hash: apiKey });
    return Array.isArray(result?.incomes) ? result.incomes : [];
  }

  async function fetchExpenses() {
    const result = await postJson('/checkExpens', { api_hash: apiKey });
    return Array.isArray(result?.expenses) ? result.expenses : [];
  }

  async function loadDashboard() {
    try {
      const [balance, incomes, expenses] = await Promise.all([
        fetchBalance(),
        fetchIncomes(),
        fetchExpenses(),
      ]);

      updateDashboard(balance, incomes, expenses);
    } catch (err) {
      console.error('Не удалось загрузить данные', err);
      alert('Не удалось загрузить данные. Попробуйте перезагрузить страницу.');
    }
  }

  function updateDashboard(balance, incomes, expenses) {
    const incomeTotal = incomes.reduce((s, item) => s + Number(item.amount || 0), 0);
    const expenseTotal = expenses.reduce((s, item) => s + Number(item.amount || 0), 0);

    const balanceEl = document.getElementById('balance');
    const incomeEl = document.getElementById('income');
    const expenseEl = document.getElementById('expense');
    const chartTotal = document.getElementById('chartTotal');

    if (balanceEl) balanceEl.textContent = formatMoney(balance);
    if (incomeEl) incomeEl.textContent = formatMoney(incomeTotal);
    if (expenseEl) expenseEl.textContent = formatMoney(expenseTotal);
    if (chartTotal) chartTotal.textContent = formatMoney(balance);

    renderTransactions(incomes, expenses);
  }

  function renderTransactions(incomes, expenses) {
    const list = document.getElementById('transactions');
    if (!list) return;

    list.innerHTML = '';

    const items = [
      ...incomes.map(item => ({ ...item, type: 'income' })),
      ...expenses.map(item => ({ ...item, type: 'expense' })),
    ];

    items.reverse();

    // Показываем только 10 последних транзакций
    const recent = items.slice(0, 10);

    recent.forEach(item => {
      const row = document.createElement('div');
      row.className = 'transaction';

      const info = document.createElement('div');
      info.className = 'transaction-info';

      const name = document.createElement('div');
      name.className = 'transaction-name';
      name.textContent = item.description || (item.type === 'income' ? 'Доход' : 'Расход');

      const date = document.createElement('div');
      date.className = 'transaction-date';
      date.textContent = new Date().toLocaleString();

      info.appendChild(name);
      info.appendChild(date);

      const amountEl = document.createElement('div');
      amountEl.className = 'amount ' + (item.type === 'income' ? 'plus' : 'minus');
      const prefix = item.type === 'income' ? '+' : '-';
      amountEl.textContent = prefix + formatMoney(item.amount || 0);

      row.appendChild(info);
      row.appendChild(amountEl);

      list.appendChild(row);
    });
  }

  async function addTransaction(type) {
    const modal = document.getElementById('transactionModal');
    const title = document.getElementById('modalTitle');
    const submit = document.getElementById('modalSubmit');
    const amount = document.getElementById('modalAmount');
    const description = document.getElementById('modalDesc');

    title.textContent = type === 'income' ? 'Добавить доход' : 'Добавить расход';
    amount.value = '';
    description.value = '';

    submit.onclick = async () => {
      const value = Number(amount.value);
      const desc = description.value.trim();

      if (!value || value <= 0) {
        alert('Введите корректную сумму');
        return;
      }

      try {
        await postJson(type === 'income' ? '/addIncome' : '/addExpens', {
          api_hash: apiKey,
          amount: value,
          description: desc,
        });

        closeModal();
        loadDashboard();
      } catch (err) {
        console.error('Не удалось сохранить транзакцию', err);
        alert('Не удалось сохранить транзакцию');
      }
    };

    modal.classList.remove('hidden');
  }

  function closeModal() {
    const modal = document.getElementById('transactionModal');
    modal.classList.add('hidden');
  }