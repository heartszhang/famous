using famousfront.messages;
using GalaSoft.MvvmLight.Messaging;
using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Diagnostics;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Input;
using System.Windows.Interactivity;
using System.Windows.Media;
using famousfront.utils;

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
          var sc = VisualTreeExtensions.FindVisualChild<ScrollViewer>(AssociatedObject) as ScrollViewer;
            if (sc == null)
                return;
            sc.ScrollToTop();
        }
    }
    class ListBoxPreventRequestBringIntoViewBehavior : Behavior<ListBox>
    {
      protected override void OnAttached()
      {
        base.OnAttached();
        TryHook();
        AssociatedObject.Loaded += new RoutedEventHandler(OnAssociatedObject).MakeWeakSpecial(eh=>AssociatedObject.Loaded -= eh);
      }
      bool _hooked;
      void OnAssociatedObject(object sender, RoutedEventArgs e)
      {
        if (_hooked)
          return;
        TryHook();
      }
      void TryHook()
      {
        var scp = VisualTreeExtensions.FindVisualChild<ScrollContentPresenter>(AssociatedObject);
        if (scp == null)
          return;
        _hooked = true;
        scp.RequestBringIntoView += new RequestBringIntoViewEventHandler(PreventRequestBringIntoView).MakeWeakSpecial(eh =>
        {
          scp.RequestBringIntoView -= eh;
        });
      }
      void PreventRequestBringIntoView(object sender, RequestBringIntoViewEventArgs e)
      {
        e.Handled = true;
      }
    }
    class ScrollViewerPreventRequestBringIntoViewBehavior : Behavior<ScrollViewer>
    {

    }
  }
