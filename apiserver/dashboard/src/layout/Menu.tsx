import * as React from 'react';
import { useState } from 'react';
import Box from '@mui/material/Box';
import LabelIcon from '@mui/icons-material/Label';
import InvoiceIcon from '@mui/icons-material/LibraryBooks';
import OrderIcon from '@mui/icons-material/AttachMoney';
import ProductIcon from '@mui/icons-material/Collections';
import CategoryIcon from '@mui/icons-material/Bookmark';
import VisitorIcon from '@mui/icons-material/People';

import {
    useTranslate,
    DashboardMenuItem,
    MenuItemLink,
    MenuProps,
    useSidebarState,
} from 'react-admin';

/* import visitors from '../visitors';
import orders from '../orders';
import invoices from '../invoices';
import products from '../products';
import categories from '../categories';
import reviews from '../reviews'; */
import SubMenu from './SubMenu';

type MenuName = 'menuCatalog' | 'menuSales' | 'menuCustomers';

const Menu = ({ dense = false }: MenuProps) => {
    const [state, setState] = useState({
        menuCatalog: true,
        menuSales: true,
        menuCustomers: true,
    });
    const translate = useTranslate();
    const [open] = useSidebarState();

    const handleToggle = (menu: MenuName) => {
        setState(state => ({ ...state, [menu]: !state[menu] }));
    };

    return (
        <Box
            sx={{
                width: open ? 200 : 50,
                marginTop: 1,
                marginBottom: 1,
                transition: theme =>
                    theme.transitions.create('width', {
                        easing: theme.transitions.easing.sharp,
                        duration: theme.transitions.duration.leavingScreen,
                    }),
            }}
        >
            <DashboardMenuItem />
            <SubMenu
                handleToggle={() => handleToggle('menuSales')}
                isOpen={state.menuSales}
                name="pos.menu.sales"
                icon={<OrderIcon />}
                dense={dense}
            >
                <MenuItemLink
                    to="/commands"
                    state={{ _scrollToTop: true }}
                    primaryText={translate(`resources.commands.name`, {
                        smart_count: 2,
                    })}
                    leftIcon={<OrderIcon />}
                    dense={dense}
                />
                <MenuItemLink
                    to="/invoices"
                    state={{ _scrollToTop: true }}
                    primaryText={translate(`resources.invoices.name`, {
                        smart_count: 2,
                    })}
                    leftIcon={<InvoiceIcon />}
                    dense={dense}
                />
            </SubMenu>
            <SubMenu
                handleToggle={() => handleToggle('menuCatalog')}
                isOpen={state.menuCatalog}
                name="pos.menu.catalog"
                icon={<ProductIcon />}
                dense={dense}
            >
                <MenuItemLink
                    to="/products"
                    state={{ _scrollToTop: true }}
                    primaryText={translate(`resources.products.name`, {
                        smart_count: 2,
                    })}
                    leftIcon={<ProductIcon />}
                    dense={dense}
                />
                <MenuItemLink
                    to="/categories"
                    state={{ _scrollToTop: true }}
                    primaryText={translate(`resources.categories.name`, {
                        smart_count: 2,
                    })}
                    leftIcon={<CategoryIcon />}
                    dense={dense}
                />
            </SubMenu>
            <SubMenu
                handleToggle={() => handleToggle('menuCustomers')}
                isOpen={state.menuCustomers}
                name="pos.menu.customers"
                icon={<VisitorIcon />}
                dense={dense}
            >
                <MenuItemLink
                    to="/customers"
                    state={{ _scrollToTop: true }}
                    primaryText={translate(`resources.customers.name`, {
                        smart_count: 2,
                    })}
                    leftIcon={<VisitorIcon />}
                    dense={dense}
                />
                <MenuItemLink
                    to="/segments"
                    state={{ _scrollToTop: true }}
                    primaryText={translate(`resources.segments.name`, {
                        smart_count: 2,
                    })}
                    leftIcon={<LabelIcon />}
                    dense={dense}
                />
            </SubMenu>
            <MenuItemLink
                to="/reviews"
                state={{ _scrollToTop: true }}
                primaryText={translate(`resources.reviews.name`, {
                    smart_count: 2,
                })}
                //leftIcon={<reviews.icon />}
                dense={dense}
            />
        </Box>
    );
};

export default Menu;
