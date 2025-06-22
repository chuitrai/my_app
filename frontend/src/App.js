// frontend/src/App.js

import React, { useState, useEffect } from 'react';
import './App.css'; // File CSS vừa sửa

function App() {
  const [users, setUsers] = useState([]); // Đổi tên 'items' thành 'users' cho rõ nghĩa
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Giả sử endpoint là /api/users, nếu không hãy sửa lại
    const apiUrl = process.env.REACT_APP_API_URL || 'http://127.0.0.1:8080';
    // Lấy dữ liệu từ API
    // Sử dụng fetch để lấy dữ liệu từ API
    console.log("Fetching data from:", `${apiUrl}/api/users`);
    fetch(`${apiUrl}/api/users`) // Sửa endpoint nếu cần
      .then(response => response.json())
      .then(data => {
        setUsers(data || []); // Đảm bảo data không phải null/undefined
        setLoading(false);
      })
      .catch(error => {
        console.error("Lỗi khi fetch dữ liệu:", error);
        setLoading(false);
      });
  }, []);

  return (
    <div className="App">
      <header className="App-header">
        <h1>👨‍💻 Danh sách Thành viên 👨‍💻</h1>
        {loading ? (
          <p className="loading-text">Đang tải dữ liệu từ server...</p>
        ) : (
          // Sử dụng className "user-list"
          <ul className="user-list">
            {users.map(user => (
              // Sử dụng className "user-item"
              <li key={user.id} className="user-item">
                {/* Giả sử bạn có các trường name, birthday, school */}
                <p className='info'>Họ và tên: {user.name}</p>
                <p className="info">Trường học: {user.school}</p>
                <p className="info">Ngày sinh: {user.birthday}</p>
              </li>
            ))}
          </ul>
        )}
      </header>
    </div>
  );
}

export default App;