/* eslint-disable @typescript-eslint/no-unused-vars */
import { Routes, Route, Navigate } from 'react-router-dom';
import SelectNews from './components/ui/SelectNews';

export default function App() {
  return (
    <header>
      <Routes>
        <Route path="/SelectNews" element={<SelectNews />} />
      </Routes>
    </header>
  );
}
