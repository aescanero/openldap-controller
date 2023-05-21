import React, { useEffect, useState } from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend } from 'recharts';
import { useDataProvider } from 'react-admin';

interface DataItem {
  time: number;
  legend: string[];
  value: number[];
  rate: number[];
}

const formatTimeTick = (timeInSeconds : number) => {
  const date = new Date(timeInSeconds * 1000);
  const hours = date.getHours().toString().padStart(2, '0');
  const minutes = date.getMinutes().toString().padStart(2, '0');
  const seconds = date.getSeconds().toString().padStart(2, '0');
  return `${hours}:${minutes}:${seconds}`;
};

const ChartComponent: React.FC = () => {
  const dataProvider = useDataProvider();
  const [data, setData] = useState<DataItem[]>([]);

  useEffect(() => {
    // Lógica para consultar la API y actualizar los datos
    const fetchData = async () => {
      try {
        const response = await dataProvider.getOne('api/monitor', { id: 0 });
        console.log("Data response:", response);
        const { data: initialData } = response;
        console.log("Data obtained:", initialData);
        setData(initialData);


      } catch (error) {
        console.error('Failed to obtain the data:', error);
      }
    };

    fetchData();
  }, []);

  useEffect(() => {
    const interval = setInterval(() => {
      const fetchData = async () => {
        try {
          const response = await dataProvider.getOne('api/monitor', { id: 0 });
          const { data: newData } = response;
          console.log("newData obtained:", newData);

          setData(prevData => {
            if (Array.isArray(prevData)) {
              console.log("Data length PRE:", prevData.length);
              const lastValue = prevData[prevData.length - 1].value;
              console.log("lastValue:", lastValue);
              console.log("newData.rate:", newData.rate);
              newData.rate = [];
              let ix : number = 0;
              lastValue.forEach(element => {
                console.log("element:", element);
                newData.rate.push(newData.value[ix] - element);
                ix++;
                console.log("Añadiendo newData rate en setData:", newData.rate);
              });
              console.log("Añadiendo d:", newData);
              return [...prevData, ...[newData]];
            }
            return [newData];
          });
        } catch (error) {
          console.error('Failed to obtain the data:', error);
        }
      };

      fetchData();

    },5000);
    
    return () => clearInterval(interval);
  }, []);

  //const lineKeys = data.length > 0 ? Object.keys(data[0].value) : [];
  console.log("Data length:", data.length);
  console.log("Data:", data);
  
  console.log("Data length pre fetch:", data.length);
  if (data.length >= 10) {
    console.error('data >10:', data.length);
    setData((prevData) => prevData.slice(1));
    console.error('post slice data >10:', data.length);
  }


  interface l {
    l: string[];
  }

  let l: l = { l: [] };

  if (data.length > 0) {
    if (Array.isArray(data)){
       data[0].legend.forEach(element => {
        console.log("element:", element);
        l.l.push(element);
      });
    }
    console.log("data[0].legend:", data[0].legend);
  }

  const legends = l.l;

//  const names = data.length > 0 ? Array.isArray(data) ? data[0].legend : [] : [];

  return (
    <LineChart width={600} height={300} data={data}>
      <CartesianGrid strokeDasharray="3 3" />
      <XAxis
        dataKey="time" 
        tickFormatter={formatTimeTick}
      />
      <YAxis />
      <Tooltip />
      <Legend />
      {legends.map((legend, index) => (
        <Line
          key={legend}
          type="monotone"
          dataKey={`rate[${index}]`}
          stroke={`#${Math.floor(Math.random() * 16777215).toString(16)}`}
        />
      ))}
    </LineChart>
  );
};

export default ChartComponent;
