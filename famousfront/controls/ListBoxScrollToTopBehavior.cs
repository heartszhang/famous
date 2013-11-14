using famousfront.core;
using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Interactivity;
using System.Windows.Media;

namespace famousfront.controls
{
    class ListBoxScrollToTopBehavior : Behavior<ListBox>
    {
        protected override void OnAttached()
        {            
            base.OnAttached();
            var dp = DependencyPropertyDescriptor.FromProperty(ItemsControl.ItemsSourceProperty, typeof(ItemsControl));
            if (dp == null)
                return;
            dp.AddValueChanged(AssociatedObject, (new EventHandler(OnItemsSourceChanged)).MakeWeakSpecial(eh => 
            {
                dp.RemoveValueChanged(AssociatedObject, eh);
            }));
        }
        void OnItemsSourceChanged(object o, EventArgs args)
        {
            var sc = GetVisualChild<ScrollViewer>(AssociatedObject) as ScrollViewer;
            if (sc == null)
                return;
            sc.ScrollToTop();
        }
        private T GetVisualChild<T>(DependencyObject parent) where T : Visual
        {
            T child = default(T);
            int numVisuals = VisualTreeHelper.GetChildrenCount(parent);
            for (int i = 0; i < numVisuals; i++)
            {
                Visual v = (Visual)VisualTreeHelper.GetChild(parent, i);
                child = v as T;
                if (child == null)
                {
                    child = GetVisualChild<T>(v);
                }
                if (child != null)
                {
                    break;
                }
            }
            return child;
        }
    }
}
