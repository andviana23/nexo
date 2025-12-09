import { CategoriesList } from '@/components/categories/categories-list';
import { Metadata } from 'next';

export const metadata: Metadata = {
  title: 'Categorias de Serviço | NEXO',
  description: 'Gerencie as categorias de serviço da sua barbearia',
};

export default function CategoriesPage() {
  return (
    <div className="space-y-6">
      <CategoriesList />
    </div>
  );
}
