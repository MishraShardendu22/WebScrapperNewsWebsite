/* eslint-disable @typescript-eslint/no-unused-vars */
import { Routes, Route } from 'react-router-dom';
import SelectNews from './components/SelectNews';
import Home from './components/Home';

export default function App() {
  return (
    <div>
      <Routes>
        <Route path="/Home" element={<Home />} />
        <Route path="/SelectedNews" element={<SelectNews />} />
        <Route path="*" element={<Home />} />
      </Routes>
    </div>
  );
}
