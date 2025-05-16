import { BrowserRouter, Route, Routes } from 'react-router';
import './App.css';
import GamePage from './pages/Game';
import Create from './pages/Create';
import List from './pages/List';
import Layout from './layout';

function App() {
  return <BrowserRouter>
    <Routes>
      <Route element={<Layout />}>
        <Route path="/" element={<List />} />
        <Route path="/play/:gameId" element={<GamePage />} />
        <Route path="/create" element={<Create />} />
      </Route>
    </Routes>
  </BrowserRouter>
}

export default App;
