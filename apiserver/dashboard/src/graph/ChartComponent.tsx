import React, { useEffect, useState } from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend } from 'recharts';

interface DataItem {
  date: string;
  value: number;
}

const ChartComponent: React.FC = () => {
  const [data, setData] = useState<DataItem[]>([]);

  useEffect(() => {
    // Lógica para consultar la API y actualizar los datos
    const fetchData = async () => {
      try {
        const response = await fetch('/monitor');
        const data: DataItem[] = await response.json();
        setData(data);
      } catch (error) {
        console.error('Failed to obtain the data:', error);
      }
    };

    fetchData();
  }, []);

  return (
    <LineChart width={600} height={300} data={data}>
      <XAxis dataKey="date" />
      <YAxis />
      <CartesianGrid stroke="#ccc" />
      <Tooltip />
      <Legend />
      <Line type="monotone" dataKey="value" stroke="#8884d8" />
    </LineChart>
  );
};

export default ChartComponent;
