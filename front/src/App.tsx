import { BrowserRouter, Route, Routes } from 'react-router';
import './App.css';
import GamePage from './pages/Game';
import Create from './pages/Create';

function App() {
  return <BrowserRouter>
    <Routes>
      <Route path="/play/:gameId" element={<GamePage />} />
      <Route path="/create" element={<Create />} />
    </Routes>
  </BrowserRouter>
}

export default App;
