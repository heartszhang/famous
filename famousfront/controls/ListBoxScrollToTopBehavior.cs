using System;
using System.ComponentModel;
using System.Windows.Controls;
using System.Windows.Interactivity;
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
      dp.AddValueChanged(AssociatedObject, (new EventHandler(OnItemsSourceChanged)).MakeWeakSpecial(eh => dp.RemoveValueChanged(AssociatedObject, eh)));
    }
    void OnItemsSourceChanged(object o, EventArgs args)
    {
      var sc = VisualTreeExtensions.FindVisualChild<ScrollViewer>(AssociatedObject);
      if (sc == null)
        return;
      sc.ScrollToTop();
    }
  }


}
