const apiKey = localStorage.getItem('api_hash');
        if (apiKey) {
          window.location.href = '/';
        }
        function showLogin() {
          window.location.href = '/auth';
        }
        async function register() {
    // 1. Берем значения из полей ввода
    const login = document.getElementById('regLogin').value;
    const password = document.getElementById('regPassword').value;
    
    // 2. Создаем объект с данными
    const userData = {
        Login: login,
        Password: password
    };
    
    // 3. Отправляем на сервер
    try {
        const response = await fetch('/register', {
            method: 'POST',                          // Тип запроса
            headers: {
                'Content-Type': 'application/json',  // Говорим, что шлем JSON
            },
            body: JSON.stringify(userData)           // Превращаем объект в JSON строку
        });
        
        // 4. Получаем ответ от сервера
        const result = await response.json();
        // const data = await response.json();
        
        // 2. Достаем apikey из ответа
        const apiKey = result.api_hash;

        // 3. Сохраняем в localStorage
        localStorage.setItem('api_hash', apiKey);
        localStorage.setItem('username', login);
        if (response.ok) {
            console.log('Успех!', result);
            alert('Регистрация успешна!');
            window.location.href = '/';
        } else {
            console.log('Ошибка:', result);
            alert('Ошибка: ' + (result.error || 'Неизвестная ошибка'));
        }
        
    } catch (error) {
        console.error('Ошибка соединения:', error);
        alert('Не удалось отправить данные');
    }
}