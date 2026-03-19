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

const authApiBase = '' // leave empty to use same origin

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

/* ДОХОД */

function addIncome() {

  const amount = Number(incomeAmount.value);
  const source = incomeSource.value;
  const comment = incomeComment.value;

  if (!amount) {
    alert("Введите сумму дохода");
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
  incomeAmount.value = "";
  incomeSource.value = "";
  incomeComment.value = "";
}

/* РАСХОД */

function addExpense() {

  const amount = Number(expenseAmount.value);
  const category = expenseCategory.value;
  const comment = expenseComment.value;

  if (!amount) {
    alert("Введите сумму расхода");
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
  expenseAmount.value = "";
  expenseCategory.value = "";
  expenseComment.value = "";
}

/* РЕНДЕР */

function render() {

  const incomeSum = incomes.reduce((s, i) => s + i.amount, 0);
  const expenseSum = expenses.reduce((s, e) => s + e.amount, 0);

  const balance = incomeSum - expenseSum;

  balanceEl.innerText = balance.toLocaleString() + " ₽";
  incomeTotal.innerText = incomeSum.toLocaleString() + " ₽";
  expenseTotal.innerText = expenseSum.toLocaleString() + " ₽";

  renderHistory();
}

function renderHistory() {

  historyList.innerHTML = "";

  const all = [
    ...incomes.map(i => ({ ...i, type: "income" })),
    ...expenses.map(e => ({ ...e, type: "expense" }))
  ];

  all.sort((a,b) => b.id - a.id);

  all.forEach(item => {

    const li = document.createElement("li");

    if (item.type === "income") {
      li.innerHTML = `
        <span>➕ ${item.source || ""} ${item.comment ? '— ' + item.comment : ''}</span>
        <strong style="color:green">+${item.amount} ₽</strong>
      `;
    } else {
      li.innerHTML = `
        <span>➖ ${item.category || ""} ${item.comment ? '— ' + item.comment : ''}</span>
        <strong style="color:red">-${item.amount} ₽</strong>
      `;
    }

    historyList.appendChild(li);
  });
}

render();